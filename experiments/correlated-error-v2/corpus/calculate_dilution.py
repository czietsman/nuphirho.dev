def calculate_dilution(existing_shares, option_pool_pct, investment, pre_money_valuation):
    """Calculate the post-investment ownership percentages for founders and
    investors after a new funding round including an option pool.

    Args:
        existing_shares: current founder shares outstanding
        option_pool_pct: option pool as a fraction of post-money shares
        investment: dollar amount of new investment
        pre_money_valuation: company valuation for the round

    Returns:
        dict with founder_pct, investor_pct, and pool_pct
    """
    post_money = pre_money_valuation + investment

    # Price per share based on pre-money
    price_per_share = pre_money_valuation / existing_shares

    # New shares for investor
    investor_shares = investment / price_per_share

    # Option pool shares
    total_shares_after = existing_shares + investor_shares
    pool_shares = total_shares_after * option_pool_pct

    # Final totals
    final_total = existing_shares + investor_shares + pool_shares

    return {
        "founder_pct": round(existing_shares / final_total * 100, 2),
        "investor_pct": round(investor_shares / final_total * 100, 2),
        "pool_pct": round(pool_shares / final_total * 100, 2),
        "price_per_share": round(price_per_share, 2)
    }
